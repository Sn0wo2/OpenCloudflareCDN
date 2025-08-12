import {Turnstile} from "@marsidev/react-turnstile";
import "./OpenCloudflareCDN.scss";
import i18n from "i18next";
import {useState} from "react";
import {Trans, useTranslation} from "react-i18next";
import {genRayID, getRootDomain} from "./util/cloudflare.ts";

interface AppProps {
    siteKey: string | null | undefined;
    successCallback: (token: string, rayID: string) => Promise<void>;
    rayID: string | null | undefined;
}

export function OpenCloudflareCDN({rayID, siteKey, successCallback}: AppProps) {
    const {t} = useTranslation();
    const sk = siteKey || '1x00000000000000000000AA';
    const [isVerified, setIsVerified] = useState(false);
    const domain = getRootDomain();

    const rid = rayID || genRayID()

    const handleSuccess = (token: string) => {
        setIsVerified(true);
        void successCallback(rid, token);
    };

    document.title = i18n.t('page_title');

    return (
        <>
            <div className="main-wrapper" role="main">
                <div className="main-content">
                    <h1 className="zone-name-title h1">{domain}</h1>
                    {isVerified ? (
                        <div>
                            <div id="challenge-success-text" className="h2">{t('success')}</div>
                            <div className="spacer"></div>
                            <div className="core-msg spacer">{t('waiting', {domain})}</div>
                        </div>
                    ) : (
                        <div>
                            <p className="h2 spacer-bottom">{t('verifying')}</p>
                            <Turnstile
                                siteKey={sk}
                                onSuccess={handleSuccess}
                                onError={() => console.error('Turnstile error')}
                                onExpire={() => console.warn('Turnstile expired')}
                                options={{theme: "dark"}}
                            />
                            <p className="core-msg spacer-top">{t('check_connection', {domain})}</p>
                        </div>
                    )}
                </div>
            </div>
            <div className="footer" role="contentinfo">
                <div className="footer-inner">
                    <div className="clearfix diagnostic-wrapper">
                        <div className="ray-id">{"Ray ID: "}<code>{rid}</code></div>
                    </div>
                    <div className="text-center" id="footer-text">
                        <Trans i18nKey="provided_by">
                            Performance & security by  <a
                            rel="noopener noreferrer"
                            href="https://github.com/Sn0wo2/OpenCloudflareCDN"
                            target="_blank"
                        >OpenCloudflareCDN</a>
                        </Trans>
                    </div>
                </div>
            </div>
        </>
    );
}