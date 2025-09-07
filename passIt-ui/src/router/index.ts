import { createRouter, createWebHistory } from 'vue-router'
import { checkAuth, isAuthenticated } from '@/composables/useAuth';
import HomeView from '../views/HomeView.vue';
import JobsView from '../views/JobsView.vue';
import NotFoundView from '@/views/NotFoundView.vue';
import JobView from '@/views/JobView.vue';
import AddJobView from '@/views/AddJobView.vue';
import EditJobView from '@/views/EditJobView.vue';
import LoginRedirectView from '@/views/LoginRedirectView.vue';
import LogoutView from '@/views/LogoutView.vue';

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
      path: '/logout',
      name: 'logout',
      component: LogoutView,
      meta: { public: false }
    },
    {
        path: '/jobs',
        name: 'jobs',
        component: JobsView,
        meta: { public: false }
    },
    {
      path: '/users/:id',
      name: 'job',
      component: JobView,
      meta: { public: false }
    },
    {
      path: '/jobs/add',
      name: 'add-job',
      component: AddJobView,
      meta: { public: false }
    },
    {
      path: '/jobs/edit/:id',
      name: 'edit-job',
      component: EditJobView,
      meta: { public: false }
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
      return next('/dev/null');
    }
  }
  next();
});

export default router
