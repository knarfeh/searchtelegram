import * as React from 'react';

const omitBy = require('lodash/omitBy')
const isUndefined = require('lodash/isUndefined')
const defaults = require('lodash/defaults')

export type RenderFunction = (props?: any, children?: any) => Element
export type Element = React.ReactElement<any>
export type RenderComponentType<P> = React.ComponentClass<P> | React.ClassicComponentClass<P> | Element | RenderFunction;


export const RenderComponentPropType = React.PropTypes.oneOfType([
  function(props: any, propName: string, componentName: string) {
    return isUndefined(props[propName]) || (props[propName]["prototype"] instanceof React.Component)
  },
  React.PropTypes.element,
  React.PropTypes.func,
])

export function renderComponent(component: RenderComponentType<any>, props={}, children=null) {
  let isReactComponent = (
    component["prototype"] instanceof React.Component ||
    (component["prototype"] && component["prototype"].isReactComponent) ||
    typeof component === 'function'
  )
  if (isReactComponent) {
    return React.createElement(
      component as React.ComponentClass<any>,
      props, children
    )
  } else if (React.isValidElement(component)) {
    return React.cloneElement(
      component as Element,
      omitBy(props, isUndefined), children
    );
  }
  console.warn('Invalid component', component)
  return null
}
