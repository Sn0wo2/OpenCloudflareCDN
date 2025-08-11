import {Turnstile} from "@marsidev/react-turnstile";
import "./OpenCloudflareCDN.scss"
import {useState} from "react";

interface AppProps {
    siteKey: string | undefined;
    successCallback: (token: string) => Promise<void>;
    rayID: string | undefined;
}

function generateCfRayId(): string {
    const arr = new Uint8Array(8)
    crypto.getRandomValues(arr)
    return Array.from(arr).map(b => b.toString(16).padStart(2, "0")).join("")
}

function getRootDomain() {
    const parts = window.location.hostname.split('.');
    return parts.length > 2
        ? parts.slice(-2).join('.')
        : window.location.hostname;
}


export function OpenCloudflareCDN({rayID, siteKey, successCallback}: AppProps) {
    const sk = siteKey || '1x00000000000000000000AA';
    const [isVerified, setIsVerified] = useState(false);

    const handleSuccess = (token: string) => {
        setIsVerified(true);
        void successCallback(token);
    };

    return (
        <>
            <div className="main-wrapper" role="main">
                <div className="main-content">
                    <h1 className="zone-name-title h1">{getRootDomain()}</h1>
                    {isVerified ? (
                        <div>
                            <div id="challenge-success-text" className="h2">验证成功</div>
                            <div className="spacer"></div>
                            <div className="core-msg spacer">正在等待 {getRootDomain()} 响应...</div>
                        </div>
                    ) : (
                        <div>
                            <p className="h2 spacer-bottom">正在验证您是否是真人。这可能需要几秒钟时间。</p>
                            <Turnstile
                                siteKey={sk}
                                onSuccess={handleSuccess}
                                onError={() => console.error('Turnstile error')}
                                onExpire={() => console.warn('Turnstile expired')}
                                options={{theme: "dark"}}
                            />
                            <p className="core-msg spacer-top">{`继续之前，${getRootDomain()}需要先检查您的连接的安全性。`}</p>
                            <noscript>
                                <div className="h2"><span>Enable JavaScript and cookies to continue</span></div>
                            </noscript>
                        </div>
                    )}
                </div>
            </div>
            <div className="footer" role="contentinfo">
                <div className="footer-inner">
                    <div className="clearfix diagnostic-wrapper">
                        <div className="ray-id">Ray ID: <code>{rayID || generateCfRayId()}</code></div>
                    </div>
                    <div className="text-center" id="footer-text">
                        {"性能和安全由"}
                        <a rel="noopener noreferrer"
                           href="https://www.cloudflare.com?utm_source=challenge&utm_campaign=m"
                           target="_blank">Cloudflare</a>
                        {"提供"}
                    </div>
                </div>
            </div>
        </>
    );
}