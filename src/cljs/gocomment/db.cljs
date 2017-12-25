(ns gocomment.db)

(def default-db
  {:name "re-frame"
   :is-loading true
   :errors []
   :reply {:username ""
           :email ""
           :content ""}
   :comments []})