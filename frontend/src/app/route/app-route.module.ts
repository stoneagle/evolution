import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes  } from '@angular/router';
import { PoolComponent } from '../controller/pool/pool.component';
import { PoolListItemComponent } from '../controller/pool/list-item/pool-list-item.component';
import { ClassifyComponent } from '../controller/classify/classify.component';
import { ItemComponent } from '../controller/item/item.component';
import { ShellComponent } from '../base/shell/shell.component';
import { DefaultComponent } from '../base/default/default.component';
import { SignInComponent } from '../base/sign-in/sign-in.component';

const routes: Routes = [
  { path: '', redirectTo: 'stock', pathMatch: 'full' }, 
  { path: 'login', component: SignInComponent  },
  { 
    path: 'stock', component: ShellComponent,
    children:[
      { path: '', component: DefaultComponent }, 
      { path: 'pool', component: PoolComponent  },
      { path: 'pool/:id', component: PoolListItemComponent  },
    ]
  },
  { 
    path: 'config', component: ShellComponent,
    children:[
      { path: '', component: DefaultComponent }, 
      { path: 'classify', component: ClassifyComponent },
      { path: 'item', component: ItemComponent },
    ]
  }
];

@NgModule({
  imports: [ RouterModule.forRoot(routes) ],
  exports: [ RouterModule ]
})
export class AppRouteModule { 
}
