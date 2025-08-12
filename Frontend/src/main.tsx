import {StrictMode} from 'react'
import {createRoot} from 'react-dom/client'
import {OpenCloudflareCDN} from "./OpenCloudflareCDN.tsx";
import './util/i18n.ts';

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <OpenCloudflareCDN siteKey={import.meta.env.VITE_SITE_KEY} successCallback={async (token) => {
            try {
                const response = await fetch('/v1/verify', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({turnstileToken: token}),
                });
                const data = await response.json();
                if (!data.success) {
                    console.error(data.msg);
                }
            } catch (err) {
                console.error(err);
            }
            window.location.reload();
        }} rayID={null}/>
    </StrictMode>,
)