// Copyright 2015 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package cache

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/seago/com"
	"github.com/seago/seago"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/seago/contrib/cache"
)

func Test_LedisCacher(t *testing.T) {
	Convey("Test nodb cache adapter", t, func() {
		opt := cache.Options{
			Adapter:       "nodb",
			AdapterConfig: "./tmp.db",
		}

		Convey("Basic operations", func() {
			s := seago.New()
			s.Use(cache.Cacher(opt))

			s.Get("/", func(c cache.Cache) {
				So(c.Put("uname", "unknwon", 1), ShouldBeNil)
				So(c.Put("uname2", "unknwon2", 1), ShouldBeNil)
				So(c.IsExist("uname"), ShouldBeTrue)

				So(c.Get("404"), ShouldBeNil)
				So(c.Get("uname").(string), ShouldEqual, "unknwon")

				time.Sleep(2 * time.Second)
				So(c.Get("uname"), ShouldBeNil)
				time.Sleep(1 * time.Second)
				So(c.Get("uname2"), ShouldBeNil)

				So(c.Put("uname", "unknwon", 0), ShouldBeNil)
				So(c.Delete("uname"), ShouldBeNil)
				So(c.Get("uname"), ShouldBeNil)

				So(c.Put("uname", "unknwon", 0), ShouldBeNil)
				So(c.Flush(), ShouldBeNil)
				So(c.Get("uname"), ShouldBeNil)
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/", nil)
			So(err, ShouldBeNil)
			s.ServeHTTP(resp, req)

			s.Get("/id", func(c cache.Cache) {
				So(c.Incr("404"), ShouldNotBeNil)
				So(c.Decr("404"), ShouldNotBeNil)

				So(c.Put("int", 0, 0), ShouldBeNil)
				So(c.Put("int64", int64(0), 0), ShouldBeNil)
				So(c.Put("string", "hi", 0), ShouldBeNil)

				So(c.Incr("int"), ShouldBeNil)
				So(c.Incr("int64"), ShouldBeNil)

				So(c.Decr("int"), ShouldBeNil)
				So(c.Decr("int64"), ShouldBeNil)

				So(c.Incr("string"), ShouldNotBeNil)
				So(c.Decr("string"), ShouldNotBeNil)

				So(com.StrTo(c.Get("int").(string)).MustInt(), ShouldEqual, 0)
				So(com.StrTo(c.Get("int64").(string)).MustInt64(), ShouldEqual, 0)

				So(c.Flush(), ShouldBeNil)
			})

			resp = httptest.NewRecorder()
			req, err = http.NewRequest("GET", "/id", nil)
			So(err, ShouldBeNil)
			s.ServeHTTP(resp, req)
		})
	})
}
