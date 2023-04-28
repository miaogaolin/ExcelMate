<script>
    import { outputData } from "./lib/store.js";
    $: value = "";

    outputData.subscribe((v) => {
        value = v.join("\n");
    });
</script>

<div class="result-container h-1/2 w-full bg-white">
    <textarea
        class="h-full placeholder:italic placeholder:text-slate-400 block bg-white w-full border border-slate-300 rounded-md py-2 pl-2.5 pr-3 shadow-sm focus:outline-none focus:ring-1 sm:text-sm"
        placeholder="Result"
        {value}
    />
</div>
