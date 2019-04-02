<template>
    <section class="bg-white">
        <div class="row-container">
            <div class="intro text-center">
                <img src="../assets/images/smt.png">
                <div>
                    <h3> Bridge </h3>
                    <br>
                    <br>
                    <b> With bridge</b>
                    <br>
                    <b> You can exchange your assets</b>
                    <br>
                    <b>between public and private chain</b>
                </div>
            </div>

            <div class="wallet-container">
                <div class="bridge-item-container" style="margin-top: 50px;">
                    <div class="bridge-item-header"> ERC20 Token</div>
                    <div class="token-container">
                        <img class="token-icon" src="../assets/images/token.png">
                        <div class="token-content-container">
                            <div class="token-exchange-info" style="margin-bottom: 50px">
                                <div class="token-label">
                                    <div>Public chain tokens:</div>
                                    <div>
                                        {{tokenPublicChain}}
                                    </div>
                                </div>
                                <form>
                                    <div class="token-form">
                                        <label class="form-label label slot-label" for="pvcExInput">Amount: </label>
                                        <input class="form-input" id="pvcExInput" type="number" v-model.number="toPriAmount">
                                        <button type="button" class="form-button alpha-button alpha-button-primary"
                                                @click="tokenToPrivate">Exchange to private chain
                                        </button>
                                    </div>
                                </form>
                            </div>
                            <div class="token-exchange-info">
                                <div class="token-label">
                                    <div>Private chain tokens:</div>
                                    <div>
                                        {{tokenPrivateChain}}
                                    </div>
                                </div>
                                <form>
                                    <div class="token-form">
                                        <label class="form-label label slot-label" for="pbcExInput">Amount: </label>
                                        <input class="form-input" id="pbcExInput" type="number" v-model.number="toPubAmount">
                                        <button type="button" class="form-button alpha-button alpha-button-primary"
                                                @click="tokenToPublic">Exchange to public chain
                                        </button>
                                    </div>
                                </form>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="bridge-item-container" style="margin-top: 100px;">
                    <div class="bridge-item-header"> Avatar</div>
                    <div class="token-container">
                        <img class="token-icon avatar-img" v-if="avatarPrivateChain!==undefined" :src="avatarPvcPath">
                        <img class="token-icon avatar-img" v-else-if="avatarPublicChain!==undefined" :src="avatarPbcPath">
                        <div class="token-icon" v-else></div>
                        <div class="token-content-container">
                            <div class="token-exchange-info" style="margin-bottom: 50px">
                                <div class="token-label">
                                    <div><b>Chain:</b></div>
                                    <div v-if="avatarPublicChain!==undefined">
                                        <b>Public</b>
                                    </div>
                                    <div v-else-if="avatarPrivateChain!==undefined">
                                        <b>Private</b>
                                    </div>
                                </div>
                                <!--<div class="token-label">-->
                                    <!--<div>Token id:</div>-->
                                    <!--<div v-if="avatarPublicChain!==undefined">-->
                                        <!--{{avatarPublicChain.tokenId}}-->
                                    <!--</div>-->
                                    <!--<div v-else-if="avatarPrivateChain!==undefined">-->
                                        <!--{{avatarPrivateChain.tokenId}}-->
                                    <!--</div>-->
                                <!--</div>-->
                                <div class="token-label">
                                    <div>Gene:</div>
                                    <div v-if="avatarPublicChain!==undefined">
                                        {{avatarPublicChain.gene}}
                                    </div>
                                    <div v-else-if="avatarPrivateChain!==undefined">
                                        {{avatarPrivateChain.gene}}
                                    </div>
                                </div>
                                <div class="token-label">
                                    <div>Level:</div>
                                    <div v-if="avatarPublicChain!==undefined">
                                        {{avatarPublicChain.avatarLevel}}
                                    </div>
                                    <div v-else-if="avatarPrivateChain!==undefined">
                                        {{avatarPrivateChain.avatarLevel}}
                                    </div>
                                </div>
                                <div class="token-label">
                                    <div>Weaponed:</div>
                                    <div v-if="avatarPublicChain!==undefined">
                                        {{avatarPublicChain.weaponed? 'false':'true'}}
                                    </div>
                                    <div v-else-if="avatarPrivateChain!==undefined">
                                        {{avatarPrivateChain.weaponed? 'false':'true'}}
                                    </div>
                                </div>
                            </div>
                            <div class="avatar-exchange-info" v-if="avatarPublicChain!==undefined">
                                <button type="button" class="form-button alpha-button alpha-button-primary"
                                        @click="avatarToPrivate">Exchange to private chain
                                </button>
                            </div>
                            <div class="avatar-exchange-info" v-else>
                                <button type="button" class="form-button alpha-button alpha-button-primary"
                                        @click="avatarToPublic">Exchange to public chain
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
    import {
        encodePrivateChainFunction, encodePublicChainFunction,
        getPrivateChainBaseTxObject,
        getPublicChainBaseTxObject,
        signTx
    } from "../assets/js/tx";

    export default {
        props: {
            avatarPrivateChain: Object,
            avatarPublicChain: Object,
            tokenPrivateChain: Number,
            tokenPublicChain: Number,
            avatarPvcPath: String,
            avatarPbcPath: String,
            ws: WebSocket,
        },
        data() {
            return {
                socketConstants: undefined,
                account: undefined,
                toPriAmount: 0,
                toPubAmount: 0,
                wallet: {
                    publicBalance: undefined,
                    balance: undefined,
                    publicEther: undefined
                },
            }
        },
        methods: {
            tokenToPublic: function () {
                if (this.toPubAmount > this.tokenPrivateChain) {
                    this.$store.state.notifyError('Tokens not enough')
                    return
                }
                if (this.toPubAmount <= 0) {
                    this.$store.state.notifyInfo('Require token amount >0')
                    return
                }
                this.$store.state.notifyInfo('Request has been received')
                let p1 = this.axios.get(`${this.$store.state.config.httpPath}/${this.socketConstants.chainType.private}/${this.account.address}/nonce`);
                let p2 = this.axios.get(`${this.$store.state.config.httpPath}/${this.socketConstants.chainType.private}/chainId`);
                Promise.all([p1, p2]).then(([r1, r2]) => {
                    let nonce = r1.data;
                    let chainId = r2.data;
                    let data = encodePrivateChainFunction("exchange", this.account.address, this.toPubAmount);
                    let tx = getPrivateChainBaseTxObject();
                    tx.from = this.account.address;
                    tx.nonce = nonce;
                    tx.data = data;
                    tx.chainId = chainId;
                    return signTx(tx, '0x' + this.account.privateKey)
                }).then((rawTx) => {
                    let payload = {
                        amount: this.toPubAmount,
                        gcuid: this.socketConstants.gcuid.exchange,
                        type: this.socketConstants.exchangeType.erc20Token,
                        source: this.socketConstants.sourceType.privateChain,
                        transaction: rawTx.slice(2),
                    };
                    this.ws.send(JSON.stringify(payload))
                }).catch(err => {
                    console.log(err)
                });
            },
            tokenToPrivate: function () {
                if (this.toPriAmount > this.tokenPublicChain) {
                    this.$store.state.notifyError('Tokens not enough')
                    return
                }
                if (this.toPriAmount <= 0) {
                    this.$store.state.notifyError('Require token amount>0')
                    return
                }
                this.$store.state.notifyInfo('Request has been received')
                let p1 = this.axios.get(`${this.$store.state.config.httpPath}/${this.socketConstants.chainType.public}/${this.account.address}/nonce`);
                let p2 = this.axios.get(`${this.$store.state.config.httpPath}/${this.socketConstants.chainType.public}/chainId`);
                Promise.all([p1, p2]).then(([r1, r2]) => {
                    let nonce = r1.data;
                    let chainId = r2.data;
                    let data = encodePublicChainFunction("exchange", this.account.address, this.toPriAmount);
                    let tx = getPublicChainBaseTxObject();
                    tx.from = this.account.address;
                    tx.nonce = nonce;
                    tx.data = data;
                    tx.chainId = chainId;
                    return signTx(tx, '0x' + this.account.privateKey)
                }).then((rawTx) => {
                    console.log("to private chain amount",this.toPriAmount);
                    let payload = {
                        amount: this.toPriAmount,
                        gcuid: this.socketConstants.gcuid.exchange,
                        type: this.socketConstants.exchangeType.erc20Token,
                        source: this.socketConstants.sourceType.publicChain,
                        transaction: rawTx.slice(2),
                    };
                    this.ws.send(JSON.stringify(payload))
                }).catch(err => {
                    console.log(err)
                });
            },
            avatarToPublic: function() {
                if(this.avatarPrivateChain === undefined) {
                    this.$store.state.notifyError('avatar not exists in private chain')
                    return
                }
                this.$store.state.notifyInfo('Request has been received')
                let p1 = this.axios.get(`${this.$store.state.config.httpPath}/${this.socketConstants.chainType.private}/${this.account.address}/nonce`);
                let p2 = this.axios.get(`${this.$store.state.config.httpPath}/${this.socketConstants.chainType.private}/chainId`);
                Promise.all([p1, p2]).then(([r1, r2]) => {
                    let nonce = r1.data;
                    let chainId = r2.data;
                    let data = encodePrivateChainFunction("exchangeNFT", this.avatarPrivateChain.tokenId);
                    let tx = getPrivateChainBaseTxObject();
                    tx.from = this.account.address;
                    tx.nonce = nonce;
                    tx.data = data;
                    tx.chainId = chainId;
                    return signTx(tx, '0x' + this.account.privateKey)
                }).then((rawTx) => {
                    let payload = {
                        gcuid: this.socketConstants.gcuid.exchange,
                        type: this.socketConstants.exchangeType.avatar,
                        source: this.socketConstants.sourceType.privateChain,
                        transaction: rawTx.slice(2),
                    };
                    this.ws.send(JSON.stringify(payload))
                }).catch(err => {
                    console.log(err)
                });
            },
            avatarToPrivate: function() {
                if(this.avatarPublicChain === undefined) {
                    this.$store.state.notifyError('avatar not exists in public chain')
                    return
                }
                this.$store.state.notifyInfo('Request has been received');
                let p1 = this.axios.get(`${this.$store.state.config.httpPath}/${this.socketConstants.chainType.public}/${this.account.address}/nonce`);
                let p2 = this.axios.get(`${this.$store.state.config.httpPath}/${this.socketConstants.chainType.public}/chainId`);
                Promise.all([p1, p2]).then(([r1, r2]) => {
                    let nonce = r1.data;
                    let chainId = r2.data;
                    let data = encodePublicChainFunction("exchangeNFT", this.avatarPublicChain.tokenId);
                    let tx = getPublicChainBaseTxObject();
                    tx.from = this.account.address;
                    tx.nonce = nonce;
                    tx.data = data;
                    tx.chainId = chainId;
                    return signTx(tx, '0x' + this.account.privateKey)
                }).then((rawTx) => {
                    let payload = {
                        gcuid: this.socketConstants.gcuid.exchange,
                        type: this.socketConstants.exchangeType.avatar,
                        source: this.socketConstants.sourceType.publicChain,
                        transaction: rawTx.slice(2),
                    };
                    this.ws.send(JSON.stringify(payload))
                }).catch(err => {
                    console.log(err)
                });
            },
        },
        created: function () {
            this.account = this.$store.state.account;
            this.socketConstants = this.$store.state.socketConstants;
        },
    }
</script>
