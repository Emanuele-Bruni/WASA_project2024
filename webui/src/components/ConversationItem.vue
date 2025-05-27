<script>
export default {
    props : ['photo','title','isGroup']
}
</script>


<template>
  <div class="conversationCard">
    <img :src="photo" class="avatar" />

    <div class="info">
      <h3>{{ title }}</h3>
    </div>

    <span class="statusDot" :class="isGroup ? 'group' : 'private'"></span>
  </div>
</template>


<style scoped>
.conversationCard {
  display: flex;
  align-items: center;
  position: relative;
  padding: 0.75rem;
  margin: 0.5rem 0;
  background-color: #f9f9f9;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 1rem;
}

.info h3 {
  margin: 0;
  font-size: 1.1rem;
}

/* Pallino in basso a destra */
.statusDot {
  position: absolute;
  bottom: 10px;
  right: 10px;
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.statusDot.group {
  background-color: red;
}

.statusDot.private {
  background-color: gray;
}
</style>
