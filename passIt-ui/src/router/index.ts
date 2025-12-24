import { createRouter, createWebHistory } from 'vue-router'
import { checkAuth, isAuthenticated } from '@/composables/useAuth';
import HomeView from '../views/HomeView.vue';
import NotFoundView from '@/views/NotFoundView.vue';
import LoginRedirectView from '@/views/LoginRedirectView.vue';
import LogoutView from '@/views/LogoutView.vue';
import SignupView from '@/views/SignupView.vue';
import ProfileView from '@/views/ProfileView.vue';
import UsersView from '@/views/UsersView.vue';
import DeletedUsersView from '@/views/DeletedUsersView.vue';
import AddUserView from '@/views/AddUserView.vue';
import EditUserView from '@/views/EditUserView.vue';
import UserProfileView from '@/views/UserProfileView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { public: true }
    },
    {
      path: '/login',
      name: 'login',
      component: LoginRedirectView,
      meta: { public: true }
    },
    {
      path: '/auth/login',
      redirect: '/login'
    },
    {
      path: '/signup',
      name: 'signup',
      component: SignupView,
      meta: { public: true }
    },
    {
      path: '/logout',
      name: 'logout',
      component: LogoutView,
      meta: { public: false }
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
      meta: { public: false }
    },
    {
      path: '/users',
      name: 'users',
      component: UsersView,
      meta: { public: false, requiresAdmin: true }
    },
    {
      path: '/users/deleted',
      name: 'deleted-users',
      component: DeletedUsersView,
      meta: { public: false, requiresAdmin: true }
    },
    {
      path: '/users/add',
      name: 'add-user',
      component: AddUserView,
      meta: { public: false, requiresAdmin: true }
    },
    {
      path: '/users/:id',
      name: 'user-profile',
      component: UserProfileView,
      meta: { public: false, requiresAdmin: true }
    },
    {
      path: '/users/edit/:id',
      name: 'edit-user',
      component: EditUserView,
      meta: { public: false, requiresAdmin: true }
    },
    {
      path: '/:catchAll(.*)',
      name: 'not-found',
      component: NotFoundView,
      meta: { public: true }
    }
    // {
    //   path: '/about',
    //   name: 'about',
    //   // route level code-splitting
    //   // this generates a separate chunk (About.[hash].js) for this route
    //   // which is lazy-loaded when the route is visited.
    //   component: () => import('../views/AboutView.vue'),
    // },
  ],
})

// GOOD
router.beforeEach(async (to, from, next) => {
  if (!to.meta.public) {
    await checkAuth();
    if (!isAuthenticated.value) {
      return next('/login');
    }
  }
  next();
});

export default router
