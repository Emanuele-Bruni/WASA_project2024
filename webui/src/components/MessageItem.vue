<script>
export default {
    props: ['messageId', 'author', 'authorId', 'text', 'hours','isMine', 'comments', 'forwardedIds'],
    methods : {
      isPhotoMessage(text) {
        return text.startsWith("[foto]");
      }
    }
}
</script>

<template>
  <div :class="['messageItem',isMine ? 'mine':'']">
    <div :class="['messageBox',isMine ? 'mineBox':'']">
      <div class="messageHeader">
        <strong>{{ author }}</strong>
        <button class="optionsButton" @click="$emit('messageOptions', messageId)">⋮</button>
      </div>
      <div v-if="forwardedIds[messageId]" class="forwardedLabel">Inoltrato</div>
      <div v-if="isPhotoMessage(text)">
        <img :src="'data:image/png;base64,' + text.slice(6)" class="messagePhoto"/>
      </div>
      <div v-else class="messageText" >{{ text }}</div>
      <div class="messageTime">{{ hours }}</div>
      <!-- Mostra i commenti se ci sono -->
      <div v-if="comments && comments.length" class="commentsList">
        <div v-for="comment in comments" :key="comment.id" class="commentItem">
          {{ comment.emoji }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.messagePhoto {
  max-width: 200px;
  max-height: 200px;
  border-radius: 0.75rem;
  margin-top: 0.5rem;
  object-fit: cover;
  border: 1px solid #ccc;
}
.forwardedLabel{
  font-style: italic;
  color: gray;
  font-size: 0.75rem;
  margin-bottom: 0.25rem;
  display: block;
}
.commentsList {
  margin-top: 0.5rem;
  padding-left: 1rem;
  border-left: 2px solid #046785;
  display: flex;
  align-items: baseline;
}

.commentItem {
  font-size: 0.85rem;
  margin-bottom: 0.25rem;
}
.messageItem {
  display: flex;
  justify-content: flex-start;
  margin: 0.5rem 0;
  width: 100%;
}
.messageItem.mine{
  justify-content: flex-end;
}
.messageBox {
  background-color: #f0f7fc;
  border: 1px solid #046785;
  border-radius: 1rem;
  padding: 1rem;
  position: relative;
  box-shadow: 1px 3px 5px rgba(0, 0, 0, 0.1);
  max-width: 70%; /* Limita la larghezza massima */
  word-wrap: break-word; /* Spezza il testo lungo */
}
.messageBox.mineBox {
  background-color: #d0f0c0; /* verde chiaro */
  border: 1px solid #5ca75c;
}
.messageHeader {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.optionsButton {
  background: none;
  border: none;
  font-size: 1.25rem;
  cursor: pointer;
  color: #046785;
  padding: 0;
  margin-left: 10px;
}

.optionsButton:hover {
  color: black;
}

.messageText {
  margin-top: 0.5rem;
  font-size: 1rem;
  font-family: "SF Pro", sans-serif;
  color: #04354c;
}

.messageTime {
  font-size: 0.75rem;
  color: #6c7a89;
  text-align: right;
  margin-top: 0.5rem;
}
</style>
