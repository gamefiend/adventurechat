package adventurechat

import "fmt"

func (c *Client) Say(Object string) string {
	fmt.Fprintf(c.connection, "you say %q", Object)
	return fmt.Sprintf("%d says '%s'", c.ID, Object)
}

func (c *Client) Go(Object string) string {
	fmt.Fprintf(c.connection, "You enter %q\n", Object)
	fmt.Fprintf(c.connection, c.room.Description)
	return fmt.Sprintf("%d goes %q\n", c.ID, Object)
}
