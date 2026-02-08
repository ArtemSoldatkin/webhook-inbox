<script lang="ts">
	import { page } from '$app/state';

	const { endpointID, webhookID, eventID } = page.params;
</script>

<nav>
	<ul>
		<li><a href={`/endpoints/${endpointID}/webhooks/${webhookID}/events`}>Events</a></li>
		<li>
			<a href={`/endpoints/${endpointID}/webhooks/${webhookID}/events/${eventID}/deliveries`}
				>Deliveries</a
			>
		</li>
	</ul>
</nav>
<h2>Events - endpoint {endpointID} webhook {webhookID} event {eventID}</h2>
<p>Events - endpoint {endpointID} webhook {webhookID} event {eventID} full description.</p>
