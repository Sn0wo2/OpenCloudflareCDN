export const genRayID = () => {
    return [...Array(16)].map(() => Math.floor(Math.random() * 16).toString(16)).join('');
}

export const getRootDomain = () => {
    const parts = window.location.hostname.split('.');
    return parts.length > 2
        ? parts.slice(-2).join('.')
        : window.location.hostname;
}