import { beforeEach, describe, expect, it, vi } from 'vitest';

const registerLanguage = vi.fn();

vi.mock('highlight.js/lib/core', () => ({
	default: {
		registerLanguage
	}
}));

vi.mock('highlight.js/lib/languages/json', () => ({
	default: 'json-language'
}));

vi.mock('highlight.js/lib/languages/xml', () => ({
	default: 'xml-language'
}));

describe('highlight-init', () => {
	beforeEach(() => {
		registerLanguage.mockClear();
		vi.resetModules();
	});

	it('registers json and xml highlighters on import', async () => {
		await import('./highlight-init');

		expect(registerLanguage).toHaveBeenCalledTimes(2);
		expect(registerLanguage).toHaveBeenNthCalledWith(1, 'json', 'json-language');
		expect(registerLanguage).toHaveBeenNthCalledWith(2, 'xml', 'xml-language');
	});
});
