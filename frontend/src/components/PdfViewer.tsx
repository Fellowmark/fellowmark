import { LinearProgress } from "@material-ui/core";
import { FC, useState } from "react";
import {
  PdfLoader,
  PdfHighlighter,
  Tip,
  Highlight,
  Popup,
  AreaHighlight,
  NewHighlight,
  IHighlight,
} from "react-pdf-highlighter";

interface AnnotatorProps {
  url: string;
  highlights: Array<IHighlight>;
  setHighlights: (highlights: Array<IHighlight>) => void;
}

const HighlightPopup = ({
  comment,
}: {
  comment: { text: string; emoji: string };
}) =>
  comment.text ? (
    <div className="Highlight__popup">
      {comment.emoji} {comment.text}
    </div>
  ) : null;

export const Annotator: FC<AnnotatorProps> = ({
  url,
  setHighlights,
  highlights,
}) => {
  const [scrollViewerTo, setScrollViewerTo] =
    useState<(highlight: IHighlight) => void>(null);
  const resetHash = () => {
    document.location.hash = "";
  };

  const parseIdFromHash = () =>
    document.location.hash.slice("#highlight-".length);

  const scrollToHighlightFromHash = () => {
    const highlight = getHighlightById(parseIdFromHash());

    if (highlight) {
      scrollViewerTo(highlight);
    }
  };

  const getHighlightById = (id: string) => {
    return highlights.find((highlight) => highlight.id === id);
  };

  const getNextId = () => String(Math.random()).slice(2);

  const addHighlight = (highlight: NewHighlight) => {
    setHighlights([{ ...highlight, id: getNextId() }, ...highlights]);
  };

  const updateHighlight = (
    highlightId: string,
    position: Object,
    content: Object
  ) => {
    console.log("Updating highlight", highlightId, position, content);

    setHighlights(
      highlights.map((h) => {
        const {
          id,
          position: originalPosition,
          content: originalContent,
          ...rest
        } = h;
        return id === highlightId
          ? {
              id,
              position: { ...originalPosition, ...position },
              content: { ...originalContent, ...content },
              ...rest,
            }
          : h;
      })
    );
  };

  return (
    <PdfLoader
      url={url}
      beforeLoad={<LinearProgress className="progressBar" />}
    >
      {(pdfDocument) => (
        <PdfHighlighter
          pdfDocument={pdfDocument}
          enableAreaSelection={(event) => event.altKey}
          onScrollChange={resetHash}
          // pdfScaleValue="page-width"
          scrollRef={(scrollTo) => {
            setScrollViewerTo(scrollTo);

            scrollToHighlightFromHash();
          }}
          onSelectionFinished={(
            position,
            content,
            hideTipAndSelection,
            transformSelection
          ) => (
            <Tip
              onOpen={transformSelection}
              onConfirm={(comment) => {
                addHighlight({ content, position, comment });

                hideTipAndSelection();
              }}
            />
          )}
          highlightTransform={(
            highlight,
            index,
            setTip,
            hideTip,
            viewportToScaled,
            screenshot,
            isScrolledTo
          ) => {
            const isTextHighlight = !Boolean(
              highlight.content && highlight.content.image
            );

            const component = isTextHighlight ? (
              <Highlight
                isScrolledTo={isScrolledTo}
                position={highlight.position}
                comment={highlight.comment}
              />
            ) : (
              <AreaHighlight
                isScrolledTo={isScrolledTo}
                highlight={highlight}
                onChange={(boundingRect) => {
                  updateHighlight(
                    highlight.id,
                    { boundingRect: viewportToScaled(boundingRect) },
                    { image: screenshot(boundingRect) }
                  );
                }}
              />
            );

            return (
              <Popup
                popupContent={<HighlightPopup {...highlight} />}
                onMouseOver={(popupContent) =>
                  setTip(highlight, (highlight) => popupContent)
                }
                onMouseOut={hideTip}
                key={index}
                children={component}
              />
            );
          }}
          highlights={highlights}
        />
      )}
    </PdfLoader>
  );
};
