import hljs from 'highlight.js/lib/core';
import json from 'highlight.js/lib/languages/json';
import xml from 'highlight.js/lib/languages/xml';

/** Registers the syntax highlighters used by body viewers. */
hljs.registerLanguage('json', json);
hljs.registerLanguage('xml', xml);
