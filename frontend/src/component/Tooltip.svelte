<script>
    export let visible = false;
    export let text, left, top;
</script>

<div
    class:hidden={!visible}
    class="text-sm absolute tooltip z-20 dark:bg-gray-500 bg-gray-400 rounded text-white p-2"
    style="left:{left};top:{top}"
>
    {text}
</div>

<style>
    @keyframes fadeInOut {
        0% {
            opacity: 0;
        }

        50% {
            opacity: 1;
        }
    }
    .tooltip {
        animation-name: fadeInOut;
        animation-duration: 1.5s;
    }
</style>
