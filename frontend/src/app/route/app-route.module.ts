import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes  } from '@angular/router';
import { PoolComponent } from '../controller/pool/pool.component';
import { ClassifyComponent } from '../controller/classify/classify.component';
import { ItemComponent } from '../controller/item/item.component';
import { ShellComponent } from '../base/shell/shell.component';
import { SignInComponent } from '../base/sign-in/sign-in.component';

const routes: Routes = [
  { path: '', redirectTo: 'login', pathMatch: 'full' }, 
  { path: 'login', component: SignInComponent  },
  { path: 'quant', component: ShellComponent, 
      children:[
        { path: 'pool', component: PoolComponent  },
        { path: 'classify', component: ClassifyComponent },
        { path: 'item', component: ItemComponent },
      ]
  },
];

@NgModule({
  imports: [ RouterModule.forRoot(routes) ],
  exports: [ RouterModule ]
})
export class AppRouteModule { 
}
