import { request } from "../index"

export const getTask = () => request.post("/get-task")
export const submit = (req) => request.post("/submit", req)

