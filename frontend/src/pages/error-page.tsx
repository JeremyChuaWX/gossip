import { useRouteError, isRouteErrorResponse } from "react-router-dom";

function ErrorPage() {
  const error = useRouteError();

  if (isRouteErrorResponse(error)) {
    return <div>{error.statusText}</div>;
  } else {
    return <div>An error occured</div>;
  }
}

export default ErrorPage;
