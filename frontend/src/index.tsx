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
import HomePage from "./pages/home-page";
import PostPage from "./pages/post-page";
import ErrorPage from "./pages/error-page";
import UserPage from "./pages/user-page";
import SignInPage from "./pages/signin-page";
import SignUpPage from "./pages/signup-page";

// layouts
import BaseLayout from "./layouts/base-layout";

// loaders
import { getUserLoader } from "./utils/react-router/user";

const queryClient = new QueryClient();

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path="/" element={<BaseLayout />} errorElement={<ErrorPage />}>
      <Route index element={<HomePage />} />

      <Route path="auth">
        <Route path="signin" element={<SignInPage />} />
        <Route path="signup" element={<SignUpPage />} />
      </Route>

      <Route path="post/:id">
        <Route index element={<PostPage />} />
        <Route path="comment/:id" />
      </Route>

      <Route
        path="user/:id"
        element={<UserPage />}
        loader={getUserLoader(queryClient)}
      />

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
