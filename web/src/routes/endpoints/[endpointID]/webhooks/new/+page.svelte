<script lang="ts">
	import { page } from '$app/state';

	const { endpointID } = page.params;
</script>

<nav>
	<ul>
		<li><a href={`/endpoints/${endpointID}/webhooks`}>Webhooks</a></li>
	</ul>
</nav>
<h2>Register webhook for endpoint {endpointID}</h2>
