import { compile } from 'svelte/compiler';
import { describe, expect, it } from 'vitest';
import pageSource from './+page.svelte?raw';

describe('routes/sources/[sourceID]/[eventID]/+page.svelte', () => {
	it('compiles successfully as a Svelte component', () => {
		expect(() =>
			compile(pageSource, {
				filename: 'src/routes/sources/[sourceID]/[eventID]/+page.svelte',
				generate: 'client'
			})
		).not.toThrow();
	});

	it('shows a missing source or event alert and otherwise renders EventView', () => {
		expect(pageSource).toContain('No source/event ID provided.');
		expect(pageSource).toContain('<EventView {sourceID} {eventID} />');
		expect(pageSource).toContain('{#if eventID && sourceID}');
	});

	it('builds breadcrumb and heading content from the route source and event ids', () => {
		expect(pageSource).toContain("label: 'Sources'");
		expect(pageSource).toContain("label: sourceID || 'Unknown source'");
		expect(pageSource).toContain("label: eventID || 'Unknown event'");
		expect(pageSource).toContain('title={`Event ID: ${eventID}`}');
		expect(pageSource).toContain('href: resolve(`/sources/${sourceID}`)');
	});
});
