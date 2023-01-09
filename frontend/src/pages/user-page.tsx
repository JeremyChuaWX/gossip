import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { getUserQuery } from "../utils/react-query/user";

function UserPage() {
  const params: any = useParams();
  const { data: getUserRes } = useQuery(getUserQuery(params.id));
  const user = getUserRes?.data;

  return <div>hello {user?.username}</div>;
}

export default UserPage;
