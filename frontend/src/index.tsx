import { StrictMode } from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import reportWebVitals from "./reportWebVitals";
import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
  RouterProvider,
} from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

// pages
import HomePage, { homePageLoader } from "./pages/home-page";
import PostPage, { postPageLoader } from "./pages/post-page";
import ErrorPage from "./pages/error-page";
import UserPage, { userPageLoader } from "./pages/user-page";
import ProfilePage, { profilePageLoader } from "./pages/profile-page";
import SignInPage from "./pages/signin-page";
import SignUpPage from "./pages/signup-page";

// layouts
import BaseLayout, { baseLayoutLoader } from "./layouts/base-layout";

const TWO_MINS = 1000 * 60 * 2;

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: TWO_MINS,
    },
  },
});

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route
      path="/"
      element={<BaseLayout />}
      errorElement={<ErrorPage />}
      loader={baseLayoutLoader(queryClient)}
    >
      <Route
        index
        element={<HomePage />}
        loader={homePageLoader(queryClient)}
      />

      <Route path="auth">
        <Route path="signin" element={<SignInPage />} />
        <Route path="signup" element={<SignUpPage />} />
      </Route>

      <Route
        path="profile"
        element={<ProfilePage />}
        loader={profilePageLoader(queryClient)}
      />

      <Route
        path="user/:id"
        element={<UserPage />}
        loader={userPageLoader(queryClient)}
      />

      <Route path="post/:id">
        <Route
          index
          element={<PostPage />}
          loader={postPageLoader(queryClient)}
        />
        <Route path="comment/:id" />
      </Route>

      <Route path="*" element={<ErrorPage />} />
    </Route>
  )
);

const root = ReactDOM.createRoot(
  document.getElementById("root") as HTMLElement
);

root.render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
