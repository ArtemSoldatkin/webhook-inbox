<script lang="ts">
	import { page } from '$app/state';

	const { endpointID, webhookID } = page.params;
</script>

<nav>
	<ul>
		<li><a href={`/endpoints/${endpointID}/webhooks/${webhookID}/events`}>Events</a></li>
	</ul>
</nav>
<h2>Webhook - endpoint {endpointID} webhook {webhookID}</h2>
<p>Webhook - endpoint {endpointID} webhook {webhookID} full description.</p>
