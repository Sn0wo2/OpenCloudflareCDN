import {Turnstile, type TurnstileInstance} from "@marsidev/react-turnstile";
import "@/OpenCloudflareCDN.scss";
import i18n from "i18next";
import {useRef, useState} from "react";
import {Trans, useTranslation} from "react-i18next";
import {genRayID, getRootDomain} from "@/util/cloudflare";

interface AppProps {
    siteKey: string | null | undefined;
    successCallback: (token: string, rayID: string) => Promise<void>;
    rayID: string | null | undefined;
}

type VerificationStatus = 'verify' | 'success' | 'error' | 'expire';

export function OpenCloudflareCDN({rayID, siteKey, successCallback}: AppProps) {
    const {t} = useTranslation();
    const [status, setStatus] = useState<VerificationStatus>('verify');
    const turnstileRef = useRef<TurnstileInstance>(null);
    const domain = getRootDomain();
    const rid = rayID || genRayID();

    document.title = i18n.t('page_title');

    return (
        <>
            <div className="main-wrapper" role="main">
                <div className="main-content">
                    <h1 className="zone-name-title h1">{domain}</h1>
                    {status === 'success' ? (
                        <div>
                            <div id="challenge-success-text" className="h2">{t('challenge_success')}</div>
                            <div className="spacer"></div>
                            <div className="core-msg spacer">{t('waiting', {domain})}</div>
                        </div>
                    ) : status === 'error' || status === 'expire' ? (
                        <div>
                            <div id="challenge-error-text" className="h2">{t(`challenge_${status}`)}</div>
                            <div className="spacer"></div>
                            <div className="core-msg spacer">{t(`challenge_${status}_description`)}</div>
                            <div className="retry-wrapper">
                                <button className="retry-button" onClick={() => {
                                    setStatus('verify');
                                    turnstileRef.current?.reset();
                                }}>{t('retry')}</button>
                            </div>
                        </div>
                    ) : (
                        <div>
                            <p className="h2 spacer-bottom">{t('verify')}</p>
                            <Turnstile
                                ref={turnstileRef}
                                siteKey={siteKey || '1x00000000000000000000AA'}
                                onSuccess={(token) => {
                                    setStatus('success');
                                    void successCallback(token, rid);
                                }}
                                onError={() => setStatus('error')}
                                onExpire={() => setStatus('expire')}
                                options={{
                                    language: i18n.language,
                                }}
                            />
                            <p className="core-msg spacer-top">{t('check_connection', {domain})}</p>
                        </div>
                    )}
                </div>
            </div>
            <div className="footer" role="contentinfo">
                <div className="footer-inner">
                    <div className="clearfix diagnostic-wrapper">
                        <div className="ray-id">Ray ID: <code>{rid}</code></div>
                    </div>
                    <div className="text-center" id="footer-text">
                        <Trans i18nKey="provided_by">
                            Performance & security by <a
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