export const updateContext = (state, action) => {
    switch (action.type) {
        case "AUTHENTICATED":
            return { ...state, user: action.payload };
        case "UNAUTHENTICATED":
            return { ...state, user: action.payload };
        case "MODULE":
            return { ...state, module: action.payload };
        default:
            return state;
    }
};


