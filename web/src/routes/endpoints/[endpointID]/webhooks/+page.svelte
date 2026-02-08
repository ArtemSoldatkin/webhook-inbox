<script lang="ts">
	import { page } from '$app/state';

	const { endpointID } = page.params;
</script>

<h2>Webhooks for endpoint {endpointID}</h2>
<nav>
	<ul>
		<li><a href={`/endpoints/${endpointID}/webhooks/new`}>Create webhook</a></li>
		<li><a href={`/endpoints/${endpointID}/webhooks/list`}>List webhooks</a></li>
	</ul>
</nav>
