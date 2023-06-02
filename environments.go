package sense

/*
// This is constructed like any other exported API but for now we'll just make
// it accessible through WithEnvironment() until there's a clear need for it.

// Environment is a Sense environment.  This type might be useful if you worked
// for Sense and wanted to use another environment.
type environment struct {
	ID     string
	Name   string
	ApiURL string
}

type environments []environment

// Get returns the environment with the given ID, or nil if not found.
func (e environments) Get(id string) *environment {
	// a map is overkill since there are few items and we will do this search rarely.
	for _, env := range e {
		if env.ID == id {
			return &env
		}
	}
	return nil
}

// getEnvironments returns a list of environments from the Sense API.
func (s *Client) getEnvironments(ctx context.Context) (envs environments, err error) {
	res, err1 := s.client.GetEnvironmentsWithResponse(ctx)
	if err := client.Ensure(err1, "getEnvironments", res, 200); err != nil {
		return nil, err
	}
	for _, e := range *res.JSON200 {
		envs = append(envs, environment{
			ID:     deref(e.Environment),
			Name:   deref(e.DisplayName),
			ApiURL: deref(e.ApiUrl),
		})
	}
	return envs, nil
}
*/
