export const genRayID = () => {
    return [...Array(16)].map(() => Math.floor(Math.random() * 16).toString(16)).join('');
}

export const getRootDomain = () => {
    const host = window.location.hostname;
    const isIPv4 = /^(\d{1,3}\.){3}\d{1,3}$/.test(host);
    const isIPv6 = /^\[?([a-fA-F0-9:]+)\]?$/.test(host);
    if (isIPv4 || isIPv6) return host;
    const parts = host.split('.');
    return parts.length > 2 ? parts.slice(-2).join('.') : host;
}
