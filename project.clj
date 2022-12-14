(defproject astro "0.1.0-SNAPSHOT"

  :description "Astro: Activity tracker"
  :url "https://github.com/joaofnds/astro"

  :dependencies [[buddy/buddy-auth "2.2.0"]
                 [buddy/buddy-core "1.6.0"]
                 [buddy/buddy-hashers "1.4.0"]
                 [buddy/buddy-sign "3.1.0"]
                 [ch.qos.logback/logback-classic "1.2.3"]
                 [cheshire "5.10.0"]
                 [cljs-ajax "0.8.0"]
                 [clojure.java-time "0.3.2"]
                 [com.cognitect/transit-clj "1.0.324"]
                 [com.datomic/datomic-free "0.9.5697" :exclusions [org.slf4j/log4j-over-slf4j org.slf4j/slf4j-nop com.google.guava/guava]]
                 [com.fasterxml.jackson.core/jackson-core "2.11.2"]
                 [com.fasterxml.jackson.core/jackson-databind "2.11.2"]
                 [com.google.guava/guava "25.1-jre"]
                 [com.walmartlabs/lacinia "0.32.0"]
                 [cprop "0.1.17"]
                 [day8.re-frame/http-fx "0.1.6"]
                 [expound "0.8.5"]
                 [funcool/struct "1.4.0"]
                 [io.rkn/conformity "0.5.1"]
                 [luminus-transit "0.1.2"]
                 [luminus-undertow "0.1.7"]
                 [luminus/ring-ttl-session "0.3.3"]
                 [markdown-clj "1.10.5"]
                 [metosin/jsonista "0.2.6"]
                 [metosin/muuntaja "0.6.7"]
                 [metosin/reitit "0.5.5"]
                 [metosin/ring-http-response "0.9.1"]
                 [mount "0.1.16"]
                 [nrepl "0.8.0"]
                 [org.clojure/clojure "1.10.1"]
                 [org.clojure/clojurescript "1.10.773" :scope "provided"]
                 [org.clojure/data.json "0.2.6"]
                 [org.clojure/tools.cli "1.0.194"]
                 [org.clojure/tools.logging "1.1.0"]
                 [org.webjars.npm/bulma "0.9.0"]
                 [org.webjars.npm/material-icons "0.3.1"]
                 [org.webjars/webjars-locator "0.40"]
                 [org.webjars/webjars-locator-jboss-vfs "0.1.0"]
                 [re-frame "1.0.0"]
                 [reagent "0.10.0"]
                 [ring-webjars "0.2.0"]
                 [ring/ring-core "1.8.1"]
                 [ring/ring-defaults "0.3.2"]
                 [selmer "1.12.28"]]

  :min-lein-version "2.0.0"

  :source-paths ["src/clj" "src/cljs" "src/cljc"]
  :test-paths ["test/clj"]
  :resource-paths ["resources" "target/cljsbuild"]
  :target-path "target/%s/"
  :main ^:skip-aot astro.core

  :plugins [[lein-cljsbuild "1.1.7"]]
  :clean-targets ^{:protect false}
  [:target-path [:cljsbuild :builds :app :compiler :output-dir] [:cljsbuild :builds :app :compiler :output-to]]
  :figwheel
  {:http-server-root "public"
   :server-logfile "log/figwheel-logfile.log"
   :nrepl-port 7002
   :css-dirs ["resources/public/css"]
   :nrepl-middleware [cider.piggieback/wrap-cljs-repl]}

  :profiles
  {:uberjar {:omit-source true
             :prep-tasks ["compile" ["cljsbuild" "once" "min"]]
             :cljsbuild {:builds
                         {:min
                          {:source-paths ["src/cljc" "src/cljs" "env/prod/cljs"]
                           :compiler
                           {:output-dir "target/cljsbuild/public/js"
                            :output-to "target/cljsbuild/public/js/app.js"
                            :source-map "target/cljsbuild/public/js/app.js.map"
                            :optimizations :advanced
                            :pretty-print false
                            :infer-externs true
                            :closure-warnings
                            {:externs-validation :off :non-standard-jsdoc :off}
                            :externs ["react/externs/react.js"]}}}}

             :aot :all
             :uberjar-name "astro.jar"
             :source-paths ["env/prod/clj"]
             :resource-paths ["env/prod/resources"]}

   :dev           [:project/dev :profiles/dev]
   :test          [:project/dev :project/test :profiles/test]

   :project/dev  {:jvm-opts ["-Dconf=dev-config.edn"]
                  :dependencies [[binaryage/devtools "1.0.2"]
                                 [cider/piggieback "0.5.0"]
                                 [doo "0.1.11"]
                                 [figwheel-sidecar "0.5.20"]
                                 [pjstadig/humane-test-output "0.10.0"]
                                 [prone "2020-01-17"]
                                 [re-frisk "1.3.4"]
                                 [ring/ring-devel "1.8.1"]
                                 [ring/ring-mock "0.4.0"]]
                  :plugins      [[com.jakemccrary/lein-test-refresh "0.24.1"]
                                 [jonase/eastwood "0.3.5"]
                                 [lein-doo "0.1.11"]
                                 [lein-figwheel "0.5.20"]]
                  :cljsbuild {:builds
                              {:app
                               {:source-paths ["src/cljs" "src/cljc" "env/dev/cljs"]
                                :figwheel {:on-jsload "astro.core/mount-components"}
                                :compiler
                                {:output-dir "target/cljsbuild/public/js/out"
                                 :closure-defines {"re_frame.trace.trace_enabled_QMARK_" true}
                                 :optimizations :none
                                 :preloads [re-frisk.preload]
                                 :output-to "target/cljsbuild/public/js/app.js"
                                 :asset-path "/js/out"
                                 :source-map true
                                 :main "astro.app"
                                 :pretty-print true}}}}

                  :doo {:build "test"}
                  :source-paths ["env/dev/clj"]
                  :resource-paths ["env/dev/resources"]
                  :repl-options {:init-ns user
                                 :timeout 120000}
                  :injections [(require 'pjstadig.humane-test-output)
                               (pjstadig.humane-test-output/activate!)]}
   :project/test {:jvm-opts ["-Dconf=test-config.edn"]
                  :resource-paths ["env/test/resources"]
                  :cljsbuild
                  {:builds
                   {:test
                    {:source-paths ["src/cljc" "src/cljs" "test/cljs"]
                     :compiler
                     {:output-to "target/test.js"
                      :main "astro.doo-runner"
                      :optimizations :whitespace
                      :pretty-print true}}}}}

   :profiles/dev {}
   :profiles/test {}})
