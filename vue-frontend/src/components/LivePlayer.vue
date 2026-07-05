<template>
  <div class="card-panel live-player">
    <div class="panel-title">直播画面</div>
    <div class="player-area">
      <video
        ref="videoRef"
        :src="videoSrc"
        autoplay
        loop
        muted
        controls
        playsinline
        class="player-video"
      ></video>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  defaultRoom: { type: String, default: '1001' },
  videoUrl: { type: String, default: '' }
})

const videoSrc = computed(() => {
  return props.videoUrl || '/demo.mp4'
})
</script>

<style scoped>
.player-area {
  width: 100%;
  max-height: 380px;
  aspect-ratio: 16 / 9;
  border-radius: 6px;
  overflow: hidden;
  background: #000;
}

.player-video {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
