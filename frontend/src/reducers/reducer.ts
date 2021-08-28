import { ContextPayload, ContextState } from "../context/context";

export enum AuthType {
  AUTHENTICATED = "AUTHENTICATED",
  UNAUTHENTICATED = "UNAUTHENTICATED",
  MODULE = "MODULE",
}

export const updateContext = (
  state: ContextPayload,
  action: ContextState
): ContextPayload => {
  switch (action.type) {
    case AuthType.AUTHENTICATED:
      return {
        ...state,
        user: action.payload.user,
        role: action.payload.role,
        module: null,
      };
    case AuthType.UNAUTHENTICATED:
      return null;
    case AuthType.MODULE:
      return { ...state, module: action.payload.module };
    default:
      return state;
  }
};
