(ns gocomment.core
  (:require [reagent.core :as reagent]
            [re-frame.core :as re-frame]
            [gocomment.events :as events]
            [gocomment.views :as views]
            [gocomment.config :as config]))


(defn dev-setup []
  (when config/debug?
    (enable-console-print!)
    (println "dev mode")))

(defn mount-root []
  (re-frame/clear-subscription-cache!)
  (reagent/render [views/main-view]
                  (.getElementById js/document "gocomment")))

(defn ^:export init []
  (re-frame/dispatch-sync [:initialize-db])
  (re-frame.core/dispatch [:fetch-comments])
  (dev-setup)
  (mount-root))
