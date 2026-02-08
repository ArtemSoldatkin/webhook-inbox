<script lang="ts">
	import { page } from '$app/state';

	const { endpointID } = page.params;
</script>

<h2>Webhooks for endpoint {endpointID}</h2>
<nav>
	<ul>
		<li><a href={`/endpoints/${endpointID}/webhooks`}>Webhooks</a></li>
	</ul>
</nav>
<ul>
	<li>
		<section>
			<h3><a href={`/endpoints/${endpointID}/webhooks/1`}>Webhook 1</a></h3>
		</section>
	</li>
	<li>
		<section>
			<h3><a href={`/endpoints/${endpointID}/webhooks/2`}>Webhook 2</a></h3>
		</section>
	</li>
	<li>
		<section>
			<h3><a href={`/endpoints/${endpointID}/webhooks/3`}>Webhook 3</a></h3>
		</section>
	</li>
</ul>
