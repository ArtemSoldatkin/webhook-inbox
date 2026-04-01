import { describe, expect, it } from 'vitest';
import { compile } from 'svelte/compiler';
import pageSource from './+page.svelte?raw';

describe('routes/sources/[sourceID]/+page.svelte', () => {
	it('compiles successfully as a Svelte component', () => {
		expect(() =>
			compile(pageSource, {
				filename: 'src/routes/sources/[sourceID]/+page.svelte',
				generate: 'client'
			})
		).not.toThrow();
	});

	it('shows a missing-source alert and otherwise renders SourceView', () => {
		expect(pageSource).toContain('No source ID provided.');
		expect(pageSource).toContain('<SourceView {sourceID} />');
		expect(pageSource).toContain('{#if !sourceID}');
	});

	it('builds breadcrumb and heading content from the route source id', () => {
		expect(pageSource).toContain("label: 'Sources'");
		expect(pageSource).toContain("label: sourceID || 'Unknown source'");
		expect(pageSource).toContain('title={`Source ID: ${sourceID}`}');
	});
});
