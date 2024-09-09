function getDescription() {
  const details = document.querySelector("rm-task-details");
  return details?.description || "";
}
