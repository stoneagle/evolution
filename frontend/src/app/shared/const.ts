export const enum AlertType {
  DANGER, 
  WARNING, 
  INFO, 
  SUCCESS
}

export const PageSize = {
  Small: 5,
  Normal: 10,
  Big: 20,
  Larget: 50,
}

export const FilterType = {
  Radio: "filter-radio",
  Checkbox: "filter-checkbox",
} 

export const dismissInterval = 10 * 1000;
export const httpStatusCode = {
  "Unauthorized": 401,
}

export const AuthType = {
  BasicAuth: "basic-auth",
  BasicAuthJwt: "basic-auth-jwt",
}

export const WsStatus = {
  Message: 0,
  Connect: 1,
  Disconnect: 2,
  Error: 3,
}

export const CommonRoutes = {  
  SIGN_IN: "/login",
}

