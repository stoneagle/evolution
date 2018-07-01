import { NgModule }              from '@angular/core';
import { CommonModule }          from '@angular/common';
import { RouterModule, Routes  } from '@angular/router';

import { PoolComponent as QuantPool }                 from '../controller/quant/pool/pool.component';
import { PoolListItemComponent as QuantPoolListItem } from '../controller/quant/pool/list-item/pool-list-item.component';
import { ClassifyComponent as QuantClassify }         from '../controller/quant/classify/classify.component';
import { ItemComponent as QuantItem }                 from '../controller/quant/item/item.component';

import { ProjectComponent as TimeProject }            from '../controller/time/project/project.component';
import { CountryComponent as TimeCountry }            from '../controller/time/country/country.component';
import { ShellComponent }                             from '../base/shell/shell.component';
import { DefaultComponent }                           from '../base/default/default.component';
import { SignInComponent }                            from '../base/sign-in/sign-in.component';

const routes: Routes = [
  { path: '', redirectTo: 'flow/stock', pathMatch: 'full' }, 
  { path: 'login', component: SignInComponent  },
  { 
    path: 'flow', component: ShellComponent, 
    children:[
      { path: '', redirectTo: 'stock/pool', pathMatch: 'full' }, 
      { 
        path: 'stock',
        children:[
          { path: '', component: DefaultComponent }, 
          { path: 'pool', component: QuantPool },
          { path: 'pool/:id', component: QuantPoolListItem },
        ]
      },
      { 
        path: 'config',
        children:[
          { path: '', component: DefaultComponent }, 
          { path: 'classify', component: QuantClassify },
          { path: 'item', component: QuantItem },
        ]
      }
    ]
  }, 
  { 
    path: 'time', component: ShellComponent, 
    children:[
      { 
        path: 'schedule',
        children:[
          { path: '', component: DefaultComponent }, 
          { path: 'project', component: TimeProject },
        ]
      },
      { 
        path: 'config',
        children:[
          { path: '', component: DefaultComponent }, 
          { path: 'country', component: TimeCountry },
        ]
      },
    ]
  }, 
];

@NgModule({
  imports: [ RouterModule.forRoot(routes) ],
  exports: [ RouterModule ]
})
export class AppRouteModule { 
}
