<script lang="ts">
    export let data: Record<string, string[]>;
    export let title: string;
</script>

<section>
	<h4>{title}</h4>
    <ul>
        {#each Object.entries(data) as [key, value]}
            <li><strong>{key}:</strong> {value.join(", ")}</li>
        {/each}
    </ul>
</section>
