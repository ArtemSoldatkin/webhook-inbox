<script lang="ts">
	import { page } from '$app/state';

	const { endpointID, webhookID } = page.params;
</script>

<h2>Events - endpoint {endpointID} webhook {webhookID}</h2>
<nav>
	<ul>
		<li><a href={`/endpoints/${endpointID}/webhooks/${webhookID}`}>Webhook</a></li>
	</ul>
</nav>

<ul>
	<li>
		<section>
			<h3><a href={`/endpoints/${endpointID}/webhooks/${webhookID}/events/1`}>Event 1</a></h3>
		</section>
	</li>
	<li>
		<section>
			<h3><a href={`/endpoints/${endpointID}/webhooks/${webhookID}/events/1`}>Event 2</a></h3>
		</section>
	</li>
	<li>
		<section>
			<h3><a href={`/endpoints/${endpointID}/webhooks/${webhookID}/events/1`}>Event 3</a></h3>
		</section>
	</li>
</ul>
