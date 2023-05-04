<script>
    import { init } from "./lib/store.js";
    import Config from "./Config.svelte";
    import Output from "./Output.svelte";
    import Tabs from "./Tabs.svelte";

    init();
</script>

<div class="flex h-full w-full fixed">
    <Tabs />

    <div class="flex-initial h-full w-1/3 min-w-[400px]">
        <Config />

        <Output />
    </div>
</div>
