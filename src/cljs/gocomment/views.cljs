(ns gocomment.views
  (:require [re-frame.core :as re-frame]
            [gocomment.subs :as subs]
            [cljsjs.moment]))


(defn pacman-loader
  []
  [:div {:class "gocomment__loader"}
   [:div {:class "la-pacman la-dark"}
    (map (fn [key] [:div {:key key} (range 6)]))]])

(defn comment-avatar
  "comment avatar"
  [comment]
  [:div {:class "gocomment__comment__avatar"}
   [:img {:src (str "http://api.adorable.io/avatar/40/" (:username comment))}]])


(defn time-ago
  "Given Date returns how long the Date is ago compared to this instant"
  [date]
  (-> date
      js/moment
      .fromNow))

(defn comment-header
  [comment]
  [:div {:class "gocomment__comment__header"}
   [:span {:class "gocomment__comment__user"}
    (if-not (nil? (:email comment))
      [:a {:href (str "mailto:" (:email comment))} (:username comment)]
      (:username comment))]
   [:span {:class "gocomment__comment__date"}
    (time-ago (:createdAt comment))]])


(defn comment-content
  [content]
  [:div {:class "gocomment__comment__content"}
   content])


;; comment-comp is recursive and should be forward declared
(declare comment-comp)


(defn comment-footer
  "comment footer for replies and reply form"
  [comment]
  [:div {:class "gocomment__comment__footer"}
   (map comment-comp (:replies comment))])


(defn comment-body
  "comment body including header and content"
  [comment]
  [:div {:class "gocomment__comment__body"}
   (comment-header comment)
   (comment-content (:content comment))
   (comment-footer comment)])


(defn comment-comp
  "Return gocomment comment structure"
  [comment]
  [:div {:class "gocomment__comment"
         :key (:id comment)}
   (comment-avatar comment)
   (comment-body comment)])


(defn comment-list
  "Display list of comments"
  []
  (let [comments @(re-frame/subscribe [::subs/comments])]
    [:div {:class "gocomment__comment-list"}
     (map comment-comp comments)]))


(defn footer
  "Return gocomment header view"
  []
  [:div {:class "gocomment__footer"}
   [:p "Powered by gocomment."]])


(defn comment-form
  "Return gocomment form"
  []
  (let [reply @(re-frame/subscribe [::subs/comments])])
  [:div {:class "gocomment__form"}
    [:textarea {:defaultValue "" :on-change #(re-frame/dispatch [:content-change (-> % .-target .-value)])}]
    [:br]
    [:input {:placeholder "Username" :on-change #(re-frame/dispatch [:username-change (-> % .-target .-value)])}]
    [:input {:placeholder "Email" :on-change #(re-frame/dispatch [:email-change (-> % .-target .-value)])}]
    [:button {:on-click #(re-frame/dispatch [:post-comment])}]])

(defn main-view []
  (let [name (re-frame/subscribe [::subs/name])
        is-loading @(re-frame/subscribe [::subs/is-loading])]
    [:div {:class "gocomment"}
     (if is-loading
        (pacman-loader)
        [:div
         (comment-form)
         (comment-list)])
     (footer)]))

