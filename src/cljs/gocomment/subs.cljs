(ns gocomment.subs
  (:require [re-frame.core :as re-frame]))

(re-frame/reg-sub
 ::name
 (fn [db]
   (:name db)))

(re-frame/reg-sub
  ::comments
  (fn [db]
    (:comments db)))

(re-frame/reg-sub
  ::is-loading
    (fn [db]
      (:is-loading db)))

(re-frame/reg-sub
  ::reply
  (fn [db]
    (:reply db)))