export const firebaseAuth = (state, action) => {
    switch (action.type) {
        case "AUTHENTICATED":
            return {...state, user: action.payload}
        case "UNAUTHENTICATED":
            return {...state, user: action.payload}
        default:
            return state;
    }
}
