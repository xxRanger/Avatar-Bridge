<template>
        <div class="signup">
            <div class="text-center" >
                <div class="signupHead">
                    <span v-if="account!==undefined">
                    <b>Welcome</b>
                    </span>
                    <span v-else>
                        <b>Create Account</b>
                    </span>
                </div>
                <img class="signupImg" src="@/assets/images/smile.png">
            </div>
            <div class="createAccount text-center" v-if="account!==undefined">
                Your account: {{account.address}}
                <br>
                <br>
                <button class="alpha-button alpha-button-primary" type="button" @click="login">Go Ahead !</button>
            </div>
            <div v-else>
                <form>
                    <div class="signupForm form-group">
                        <button type="button" class="signupButton alpha-button alpha-button-primary" @click="create">Create</button>
                    </div>
                </form>
            </div>
        </div>
</template>
<script>
    import {createAccount} from "../assets/js/tx";
    export default {
    data() {
        return {
            account: undefined
        }
    },
    methods: {
        create() {
            console.log("create account");
            this.account = createAccount();
            this.account.address = '0x'+this.account.address;
            this.$cookies.set('account', this.account,'7d');
            this.$store.commit('setAccount', this.account);
            console.log(this.account);
        },
        login() {
            this.$router.replace('/')
        }
    },
    mounted: function() {
        if(this.$cookies.isKey("account")) {
            this.account = this.$cookies.get("account");
            this.$store.commit('setAccount', this.account);
            console.log(this.account);
        }
    },
}
</script>
<style>
.signin {
    margin-top: 20px;
}

.formcenter {
    padding-top: 100px;
}

.formpos {
    margin-left: -50px;
}

.signup {
    margin-top: 40px;
}
</style>}