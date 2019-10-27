/*
This is the backend for the branch manager called 'branma'. It analyses your feature branches and connects it with your JIRA tickets.

Copyright (C) 2019 Lars Gaubisch

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package ping

func TestPingHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/ping", nil)
    if err != nil {
        t.Fatal(err)
    }

    w := httptest.NewRecorder()
    svc, err := smis.NewService(&http.Server{}, mux.NewRouter(), logrus.New())
    if err != nil {
        t.Fatal(err)
    }
    ep := &ping{svc: svc}
    handler := http.HandlerFunc(ep.pingHandler)
    handler.ServeHTTP(w, req)

    if status := w.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expected := "pong"
    if w.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
    }
}

func TestInit(t *testing.T) {
    router := mux.NewRouter()
    srv := &http.Server{
        Handler:      router,
        Addr:         fmt.Sprintf(":%d", 30000),
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    svc := &smis.Service{
        Log:    logrus.New(),
        Router: router,
        Server: srv,
    }

    if err := Init(svc); err != nil {
        t.Fatalf("init failed: %s", err)
    }

    err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
        pathTemplate, err := route.GetPathTemplate()
        if err != nil {
            return err
        }

        if pathTemplate != "/ping" {
            t.Errorf("Expected single endpoint '/ping' but got '%s'", pathTemplate)
        }
        return nil
    })

    if err != nil {
        t.Fatalf("walk through routes failed: %s", err)
    }
}
