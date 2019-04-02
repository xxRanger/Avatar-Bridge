import Vue from 'vue'
import Router from 'vue-router'
import store from './store.js'


// const Home = async () => {
//     return import ('@/components/Home')
// }

// const Signup = async () => {
//     return import ('@/components/Signup')
// }
// const Faucet = async () => {
//     return import ('@/components/faucet')
// }
// // const Login = async () => {
// //     return import ('@/components/Login')
// // }



Vue.use(Router);

const router = new Router({
    routes: [{
        path: '/',
        name: 'Home',
        component:() =>
            import ('@/components/home'),
        meta: {
            title: 'secdev-demo',
            requireAuth: true
        }
    }, {
        path: '/login',
        name: 'Login',
        component: ()=>import ('@/components/Login'),
        meta: {
            title: 'secdev-demo',
            requireAuth: false
        }
    }]
});

router.beforeEach((to, from, next) => {
    if (to.name === 'Login') {
        if (store.state.account !== undefined) {
            next('/');
        } else {
            next();
        }
    } else {
        if (to.meta.requireAuth && store.state.account === undefined) {
            next('/login')
        } else {
            next();
        }
    }
});


export default router;