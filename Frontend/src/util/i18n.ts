import i18n from 'i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import {initReactI18next} from 'react-i18next';
import en from "../assets/locales/en.json";
import zh from "../assets/locales/zh.json";

void i18n
    .use(LanguageDetector)
    .use(initReactI18next)
    .init({
        resources: {en: {translation: en}, zh: {translation: zh}},
        fallbackLng: 'en',
        detection: {
            order: ['querystring', 'navigator', 'htmlTag'],
        },
        interpolation: {
            escapeValue: false,
        },
    });

export default i18n;