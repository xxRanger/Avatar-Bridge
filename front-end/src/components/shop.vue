<template>
    <section class="bg-white">
        <div class="row-container">
            <div class="intro text-center">
                <img src="../assets/images/shop.png">
                <div>
                    <br>
                    <h3> Shop </h3>
                    <br>
                    <b> In private chain</b>
                    <br>
                    <b> You can use tokens to </b>
                    <br>
                    <b>buy strong weapons</b>
                    <br>
                    <b>for your avatar</b>
                </div>
            </div>

            <div class="wallet-container">
                <div class=" bridge-item-container">
                    <div class="shop-avatar-info-container">
                        <div class="role-info-header">
                            Private chain avatar
                        </div>
                        <div v-if="avatarPrivateChain!==undefined">
                            <img class="avatar-img" :src="avatarPvcPath">
                        </div>
                    </div>
                    <div class="shop-token-info-container">
                        <div class="role-info-header">
                            Private chain tokens
                        </div>
                        <div>
                            {{tokenPrivateChain}}
                        </div>
                    </div>
                </div>
                <div class="bridge-item-container">
                    <div class="bridge-item-header"> Weapon</div>
                    <div class="token-container">
                        <img class="token-icon" src="../assets/images/weapon.png">
                        <div class="token-content-container">
                            <div class="token-exchange-info" style="margin-bottom: 50px">
                                <div class="token-label">
                                    <div>Token price:</div>
                                    <div>
                                        {{weaponPrice}}
                                    </div>
                                </div>
                                <div class="token-label">
                                    <div>Attack:</div>
                                    <div>
                                        5
                                    </div>
                                </div>
                            </div>
                            <div class="avatar-exchange-info">
                                <button type="button" class="form-button alpha-button alpha-button-primary"
                                        @click="buy">Buy
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>
</template>
<script>
    import {encodePrivateChainFunction, getPrivateChainBaseTxObject, signTx} from "../assets/js/tx";

    export default {
        props: {
            ws: WebSocket,
            avatarPrivateChain: Object,
            avatarPublicChain: Object,
            tokenPrivateChain: Number,
            tokenPublicChain: Number,
            avatarPvcPath: String,
            avatarPbcPath: String,
        },
        data() {
            return {
                socketConstants: undefined,
                toPriAmount: 0,
                toPubAmount: 0,
                wallet: {
                    publicBalance: undefined,
                    balance: undefined,
                    publicEther: undefined
                },
                refreshWalletInfo: undefined,
                refreshQueueInfo: undefined,
                weaponPrice: 2,
                account: undefined,
            }
        },
        methods: {
            buy: function () {
                if (this.avatarPrivateChain === undefined) {
                    this.$store.state.notifyError('avatar not exists in private chain');
                    return
                }
                if (this.avatarPrivateChain.weaponed) {
                    this.$store.state.notifyError('avatar already has weapon');
                    return
                }
                if (this.weaponPrice > this.tokenPrivateChain) {
                    this.$store.state.notifyError('Tokens not enough');
                    return
                }
                this.$store.state.notifyInfo('Request has been received');
                let payload = {
                    gcuid: this.socketConstants.gcuid.consume,
                    type: this.socketConstants.consumeType.weapon,
                    amount:this.weaponPrice,
                    address: this.account.address,
                    tokenId: this.avatarPrivateChain.tokenId,
                };
                this.ws.send(JSON.stringify(payload));
            },
        },
        created: function () {
            this.account = this.$store.state.account;
            this.socketConstants = this.$store.state.socketConstants;
        },
    }
</script>
