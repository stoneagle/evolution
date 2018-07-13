import { NgModule }              from '@angular/core';
import { CommonModule }          from '@angular/common';
import { RouterModule, Routes  } from '@angular/router';

import { PoolComponent as QuantPool }                 from '../controller/quant/pool/pool.component';
import { PoolListItemComponent as QuantPoolListItem } from '../controller/quant/pool/list-item/list-item.component';
import { ClassifyComponent as QuantClassify }         from '../controller/quant/classify/classify.component';
import { ItemComponent as QuantItem }                 from '../controller/quant/item/item.component';

import { QuestComponent as TimeQuest }               from '../controller/time/quest/quest.component';
import { ProjectComponent as TimeProject }           from '../controller/time/project/project.component';
import { TaskComponent as TimeTask }                 from '../controller/time/task/task.component';
import { ActionComponent as TimeAction }             from '../controller/time/action/action.component';
import { CountryComponent as TimeCountry }           from '../controller/time/country/country.component';
import { FieldComponent as TimeField }               from '../controller/time/field/field.component';
import { AreaComponent as TimeArea }                 from '../controller/time/area/area.component';
import { UserResourceComponent as TimeUserResource } from '../controller/time/user-resource/user-resource.component';
import { UserComponent as SystemUser }               from '../controller/system/user/user.component';
import { ShellComponent }                            from '../base/shell/shell.component';
import { DefaultComponent }                          from '../base/default/default.component';
import { SignInComponent }                           from '../base/sign-in/sign-in.component';

const routes: Routes = [
  { path: '', redirectTo: 'sign/login', pathMatch: 'full' }, 
  { 
    path: 'flow', component: ShellComponent, 
    children:[
      { 
        path: '', 
        redirectTo: 'stock/pool', 
        pathMatch: 'full' 
      }, 
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
        path: '', 
        component: DefaultComponent,
      }, 
      { 
        path: 'schedule',
        children:[
          { path: '', component: DefaultComponent }, 
          { path: 'quest', component: TimeQuest },
          { path: 'project', component: TimeProject },
          { path: 'task', component: TimeTask },
          { path: 'action', component: TimeAction },
        ]
      },
      { 
        path: 'user',
        children:[
          { path: '', component: DefaultComponent }, 
          { path: 'resource', component: TimeUserResource },
        ]
      },
      { 
        path: 'config',
        children:[
          { path: '', component: DefaultComponent }, 
          { path: 'country', component: TimeCountry },
          { path: 'field', component: TimeField },
          { path: 'area', component: TimeArea },
        ]
      },
    ]
  }, 
  { 
    path: 'sign',
    children:[
      { 
        path: 'login', component: SignInComponent },
    ]
  }, 
  { 
    path: 'system', component: ShellComponent, 
    children:[
      { 
        path: 'sign',
        children:[
          { path: 'login', component: SignInComponent },
        ]
      },
      { 
        path: 'user',
        children:[
          { path: '', component: DefaultComponent }, 
          { path: 'resource', component: SystemUser },
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
