(ns gocomment.events
  (:require [re-frame.core :as re-frame]
            [day8.re-frame.http-fx]
            [ajax.core :as ajax]
            [gocomment.db :as db]))


(re-frame/reg-event-db
 :initialize-db
 (fn  [_ _]
   db/default-db))


(re-frame/reg-event-fx                    ;; note the trailing -fx
  :fetch-comments
  (fn [{:keys [db]} _]                    ;; the first param will be "world"
    {:db   (assoc db :is-loading true)   ;; causes the twirly-waiting-dialog to show??
     :http-xhrio {:method          :get
                  :uri             "http://localhost:8080/?url="
                  :timeout         8000                                           ;; optional see API docs
                  :response-format (ajax/json-response-format {:keywords? true})  ;; IMPORTANT!: You must provide this.
                  :on-success      [:fetch-comments-success]
                  :on-failure      [:fetch-comments-fail]}}))


(defn recursive-comments
  "Takes flat vector of comments and recurses them"
  [comments {id :id :as comment}]
  (if-let [children (get comments id)]
    (assoc comment :replies (map (partial recursive-comments comments) children))
    comment))


(defn parent-comments
  "Given vector of comments returns a map of parent-id comment-list pairs"
  [comments]
  (reduce
    (fn [parents {parent-id :parentId :as comment}]
      (assoc parents parent-id
        (if-let [children (get parents parent-id)]
          (conj children comment)
          [comment])))
    {}
    comments))


(re-frame/reg-event-db
  :fetch-comments-success
  (fn [db [_ result]]
    (let [parents (parent-comments result)]
      (->> (get parents 0)
           (map (partial recursive-comments parents))
           (#(assoc db :comments % :is-loading false))))))


(re-frame/reg-event-db
  :fetch-comments-fail
  (fn [db [_ error]]
    (assoc db :errors (conj (:errors db) {:message "Could not fetch comments."} :is-loading false))))


(re-frame/reg-event-db
  :content-change
  (fn  [db [_ content]]
    (assoc-in db [:reply :content] content)))


(re-frame/reg-event-db
  :email-change
  (fn  [db [_ content]]
    (assoc-in db [:reply :email] content)))


(re-frame/reg-event-db
  :username-change
  (fn  [db [_ content]]
    (assoc-in db [:reply :username] content)))


(re-frame/reg-event-fx
  :post-comment
  (fn [{:keys [db]} _]                    ;; the first param will be "world"
    {:db   (assoc db :is-loading true)    ;; causes the twirly-waiting-dialog to show??
     :http-xhrio {:method          :post
                  :uri             "http://localhost:8080/?url="
                  :params          (:reply db)
                  :format          (ajax/json-request-format)
                  :timeout         8000                                           ;; optional see API docs
                  :response-format (ajax/json-response-format {:keywords? true})  ;; IMPORTANT!: You must provide this.
                  :on-success      [:fetch-comments]
                  :on-failure      [:post-comment-fail]}}))



(re-frame/reg-event-db
  :post-comment-fail
  (fn [db [_ error]]
    (assoc db :errors (conj (:errors db) {:message "Could not fetch comments."} :is-loading false))))
