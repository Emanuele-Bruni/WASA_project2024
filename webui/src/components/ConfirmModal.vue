<script>
export default {
    props : {
        message : {
            type: String,
            default: "Sei sicuro di voler procedere?",
        }
    }
}
</script>

<template>
    <div class="confirm-overlay">
        <div class="confirm-content">
            <h3> {{message}} </h3>
            <div class="buttons">
                <button @click="$emit('yes')"> Conferma </button>
                <button @click="$emit('no')"> Annulla </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
h3{
    font-family: "SF Pro";
    font-size: 25px;
    color: rgb(4, 103, 133);
}
.buttons button:last-child{
    background-color: rgb(226, 22, 22);
	color: white;
}
.buttons button:first-child{
    background-color: rgb(4, 103, 133);
	color: white;
}
.buttons button{
    padding: 10px 20px;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	font-family: "SF Pro";
}
.buttons{
    display: flex;
    justify-content: space-between;
    margin-top: 1 rem;
}
.confirm-content{
    background: white;
    padding: 2rem;
    border-radius: 12px;
    width: 450px;
    max-height: 90vh;
    overflow-y: auto;
    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}
.confirm-overlay{
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0,0,0,0.4);
    display: flex;
    align-items: center;
    justify-content: center;
}
</style>